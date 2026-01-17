# Next.jsのコードからuseEffect排除の設計書

## Overview

### 目的
現在のNext.jsのコードからuseEffectを、絶対にuseEffectを利用しなければいけない箇所を除いて、useEffectを排除する。処理はイベントドリブンな処理などに置き換える。

### ユーザー
- **開発者**: useEffectの排除により、コードの可読性と保守性を向上させる
- **エンドユーザー**: 既存の機能が正常に動作することを維持

### 影響
現在のシステム状態を以下のように変更する：
- 複数のページコンポーネント: Server Componentsへの移行、またはrefコールバック関数の活用
- カスタムフック: useEffectを排除し、refコールバック関数やイベントハンドラーに置き換え
- コンポーネント: useEffectを排除し、refコールバック関数やイベントハンドラーに置き換え

### Goals
- 初期データ読み込みのuseEffectを排除する
- Intersection ObserverのuseEffectを排除する
- イベントリスナーのuseEffectを排除する
- localStorage読み込みのuseEffectを排除する
- その他のuseEffectを排除する（絶対に必要な場合を除く）
- テストコードのuseEffectも排除する（絶対に必要な場合を除く）

### Non-Goals
- サーバー側（Go）のコードの変更
- 外部ライブラリのコードの変更
- 絶対にuseEffectが必要な箇所（例: クリーンアップ処理）の変更

## Architecture

### 設計方針

#### 1. Server Componentsの活用
- 初期データ読み込みが必要なページは、可能な限りServer Componentsに移行
- Server Componentsでは、サーバー側でデータを取得し、クライアントに渡す
- ただし、既存のページが'use client'を使用している場合は、refコールバック関数を使用

#### 2. refコールバック関数の活用
- 要素がマウントされた時点で処理を実行する必要がある場合は、refコールバック関数を使用
- refコールバック関数は、要素がマウントされた時点で呼び出される
- クリーンアップが必要な場合は、refコールバック関数からクリーンアップ関数を返す

#### 3. イベントハンドラーの活用
- ユーザー操作に応じた処理は、イベントハンドラーで実装
- ボタンクリック、フォーム送信などのイベントで処理をトリガー

#### 4. CSSメディアクエリの活用
- 可能な限り、CSSメディアクエリを使用して、JavaScriptでの監視を不要にする

## 詳細設計

### 1. 初期データ読み込みの置き換え

#### 1.1 対象ファイルと設計

##### `client/app/dm-posts/page.tsx`
**現状**: `useEffect(() => { loadPosts(); loadUsers(); }, [])`で初期データを読み込み

**設計**:
- Server Componentに移行するか、refコールバック関数を使用
- ただし、このページはフォーム送信などのインタラクティブな機能があるため、Client Componentとして維持
- 初期データ読み込みは、refコールバック関数を使用して、コンポーネントがマウントされた時点で実行

**実装方法**:
```typescript
const containerRef = useRef<HTMLDivElement>(null)
const hasLoadedRef = useRef(false)

const loadInitialData = async () => {
  if (hasLoadedRef.current) return
  hasLoadedRef.current = true
  await loadPosts()
  await loadUsers()
}

// refコールバック関数
const setContainerRef = (node: HTMLDivElement | null) => {
  containerRef.current = node
  if (node && !hasLoadedRef.current) {
    loadInitialData()
  }
}
```

##### `client/app/dm-users/page.tsx`
**現状**: `useEffect(() => { loadUsers(); }, [])`で初期データを読み込み

**設計**:
- 同様に、refコールバック関数を使用

##### `client/app/dm-user-posts/page.tsx`
**現状**: `useEffect(() => { loadUserPosts(); }, [])`で初期データを読み込み

**設計**:
- 同様に、refコールバック関数を使用

##### `client/app/dm_feed/[userId]/page.tsx`
**現状**: `useEffect(() => { loadInitialPosts(); }, [userId])`で初期データを読み込み

**設計**:
- userIdが変更されたときにデータを読み込む必要がある
- refコールバック関数と、userIdの変更を検知する方法を組み合わせる
- または、userIdが変更されたときに、明示的にデータを読み込むイベントハンドラーを呼び出す

**実装方法**:
```typescript
const containerRef = useRef<HTMLDivElement>(null)
const currentUserIdRef = useRef<string | null>(null)

const setContainerRef = (node: HTMLDivElement | null) => {
  containerRef.current = node
  if (node && currentUserIdRef.current !== userId) {
    currentUserIdRef.current = userId
    loadInitialPosts()
  }
}

// userIdが変更されたときの処理
if (currentUserIdRef.current !== userId) {
  currentUserIdRef.current = userId
  loadInitialPosts()
}
```

##### `client/app/dm_feed/[userId]/[postId]/page.tsx`
**現状**: `useEffect(() => { loadInitialData(); }, [userId, postId])`で初期データを読み込み

**設計**:
- 同様に、userIdとpostIdの変更を検知してデータを読み込む

### 2. Intersection Observerの置き換え

#### 2.1 対象ファイルと設計

##### `client/app/dm_feed/[userId]/page.tsx`
**現状**: `useEffect`でIntersection Observerを登録し、スクロールを検知

**設計**:
- refコールバック関数を使用して、要素がマウントされた時点でObserverを登録
- Observerのコールバック関数内で、データ読み込み処理を実行

**実装方法**:
```typescript
const loadMoreRef = useRef<HTMLDivElement>(null)

const setLoadMoreRef = (node: HTMLDivElement | null) => {
  if (node) {
    const observer = new IntersectionObserver(
      (entries) => {
        if (entries[0].isIntersecting && hasMore && !isLoadingMore) {
          loadOlderPosts()
        }
      },
      { threshold: 0.1 }
    )
    observer.observe(node)
    
    // クリーンアップ関数を返す
    return () => {
      observer.disconnect()
    }
  }
  return undefined
}
```

##### `client/app/dm_feed/[userId]/[postId]/page.tsx`
**現状**: 同様にIntersection Observerを使用

**設計**:
- 同様に、refコールバック関数を使用

##### `client/lib/hooks/use-intersection-observer.ts`
**現状**: `useEffect`でIntersection Observerを登録

**設計**:
- カスタムフックを、refコールバック関数を返す関数に変更
- または、コンポーネント内で直接refコールバック関数を使用

**実装方法**:
```typescript
function useIntersectionObserver(
  elementRef: RefObject<Element>,
  options: Args,
): IntersectionObserverEntry | undefined {
  const [entry, setEntry] = useState<IntersectionObserverEntry>()
  const observerRef = useRef<IntersectionObserver | null>(null)

  const setElementRef = (node: Element | null) => {
    if (node && !observerRef.current) {
      const observer = new IntersectionObserver(
        (entries) => {
          setEntry(entries[0])
        },
        options
      )
      observer.observe(node)
      observerRef.current = observer
    } else if (!node && observerRef.current) {
      observerRef.current.disconnect()
      observerRef.current = null
    }
  }

  // elementRef.currentにsetElementRefを設定する方法を検討
  // または、refコールバック関数を直接返す
}
```

### 3. イベントリスナーの置き換え

#### 3.1 対象ファイルと設計

##### `client/lib/hooks/use-scroll.ts`
**現状**: `useEffect`でスクロールイベントリスナーを登録

**設計**:
- refコールバック関数を使用して、要素がマウントされた時点でイベントリスナーを登録

**実装方法**:
```typescript
export default function useScroll(threshold: number) {
  const [scrolled, setScrolled] = useState(false)
  const listenerRef = useRef<(() => void) | null>(null)

  const setScrollRef = (node: HTMLElement | null) => {
    if (node && !listenerRef.current) {
      const onScroll = () => {
        setScrolled(window.pageYOffset > threshold)
      }
      window.addEventListener("scroll", onScroll)
      listenerRef.current = () => {
        window.removeEventListener("scroll", onScroll)
      }
    } else if (!node && listenerRef.current) {
      listenerRef.current()
      listenerRef.current = null
    }
  }

  return { scrolled, setScrollRef }
}
```

##### `client/lib/hooks/use-media-query.ts`
**現状**: `useEffect`でリサイズイベントリスナーを登録

**設計**:
- refコールバック関数を使用
- または、CSSメディアクエリを活用（可能な場合）

**実装方法**:
```typescript
export default function useMediaQuery() {
  const [device, setDevice] = useState<"mobile" | "tablet" | "desktop" | null>(null)
  const [dimensions, setDimensions] = useState<{ width: number; height: number } | null>(null)
  const listenerRef = useRef<(() => void) | null>(null)

  const checkDevice = () => {
    if (window.matchMedia("(max-width: 640px)").matches) {
      setDevice("mobile")
    } else if (window.matchMedia("(min-width: 641px) and (max-width: 1024px)").matches) {
      setDevice("tablet")
    } else {
      setDevice("desktop")
    }
    setDimensions({ width: window.innerWidth, height: window.innerHeight })
  }

  const setMediaQueryRef = (node: HTMLElement | null) => {
    if (node && !listenerRef.current) {
      checkDevice()
      window.addEventListener("resize", checkDevice)
      listenerRef.current = () => {
        window.removeEventListener("resize", checkDevice)
      }
    } else if (!node && listenerRef.current) {
      listenerRef.current()
      listenerRef.current = null
    }
  }

  return {
    device,
    width: dimensions?.width,
    height: dimensions?.height,
    isMobile: device === "mobile",
    isTablet: device === "tablet",
    isDesktop: device === "desktop",
    setMediaQueryRef,
  }
}
```

### 4. localStorage読み込みの置き換え

#### 4.1 対象ファイルと設計

##### `client/lib/hooks/use-local-storage.ts`
**現状**: `useEffect`でlocalStorageから値を読み込み

**設計**:
- 初期値を直接localStorageから読み込む（SSRを考慮）
- または、refコールバック関数を使用

**実装方法**:
```typescript
const useLocalStorage = <T>(
  key: string,
  initialValue: T,
): [T, (value: T) => void] => {
  // SSRを考慮して、初期値を設定
  const [storedValue, setStoredValue] = useState<T>(() => {
    if (typeof window === 'undefined') {
      return initialValue
    }
    try {
      const item = window.localStorage.getItem(key)
      return item ? JSON.parse(item) : initialValue
    } catch (error) {
      return initialValue
    }
  })

  const setValue = (value: T) => {
    setStoredValue(value)
    if (typeof window !== 'undefined') {
      window.localStorage.setItem(key, JSON.stringify(value))
    }
  }

  return [storedValue, setValue]
}
```

### 5. その他のuseEffectの置き換え

#### 5.1 対象ファイルと設計

##### `client/components/shared/counting-numbers.tsx`
**現状**: `useEffect`でrequestAnimationFrameを使用してアニメーション

**設計**:
- valueが変更されたときに、requestAnimationFrameを開始
- refコールバック関数を使用して、アニメーションを開始

**実装方法**:
```typescript
export default function CountingNumbers({
  value,
  className,
  start = 0,
  duration = 800,
}: {
  value: number
  className: string
  start?: number
  duration?: number
}) {
  const [count, setCount] = useState(start)
  const animationRef = useRef<number | null>(null)
  const previousValueRef = useRef(value)

  const startAnimation = () => {
    if (animationRef.current) {
      cancelAnimationFrame(animationRef.current)
    }

    let startTime: number | undefined
    const animateCount = (timestamp: number) => {
      if (!startTime) startTime = timestamp
      const timePassed = timestamp - startTime
      const progress = timePassed / duration
      const currentCount = easeOutQuad(progress, start, value, 1)
      if (currentCount >= value) {
        setCount(value)
        return
      }
      setCount(currentCount)
      animationRef.current = requestAnimationFrame(animateCount)
    }
    animationRef.current = requestAnimationFrame(animateCount)
  }

  // valueが変更されたときにアニメーションを開始
  if (previousValueRef.current !== value) {
    previousValueRef.current = value
    startAnimation()
  }

  return <p className={className}>{Intl.NumberFormat().format(count)}</p>
}
```

##### `client/app/dm_movie/upload/page.tsx`
**現状**: `useEffect`でUppyインスタンスを作成

**設計**:
- refコールバック関数を使用して、要素がマウントされた時点でUppyインスタンスを作成

**実装方法**:
```typescript
const uppyRef = useRef<Uppy | null>(null)
const containerRef = useRef<HTMLDivElement>(null)

const setContainerRef = (node: HTMLDivElement | null) => {
  containerRef.current = node
  if (node && !uppyRef.current) {
    const uppyInstance = createMovieUploader({
      onUploadProgress: (percent) => {
        setUploadProgress(percent)
      },
      onUploadSuccess: () => {
        setUploadStatus('success')
        setUploadProgress(100)
      },
      onUploadError: (error) => {
        setUploadStatus('error')
        setErrorMessage(error)
      },
      onUploadStart: () => {
        setUploadStatus('uploading')
        setUploadProgress(0)
        setErrorMessage(null)
      },
    })
    uppyRef.current = uppyInstance
    setUppy(uppyInstance)
    
    // クリーンアップ関数を返す
    return () => {
      uppyInstance.destroy()
      uppyRef.current = null
    }
  }
  return undefined
}
```

### 6. 保持が必要なuseEffect

#### 6.1 対象ファイルと設計

##### `client/components/video-player/video-player.tsx`
**現状**: `useEffect`でHLS.jsのクリーンアップを実行

**設計**:
- コンポーネントのアンマウント時にリソースをクリーンアップする必要があるため、useEffectを保持
- これはReactのライフサイクルに依存する処理であり、useEffectが適切

## 実装上の注意事項

### 1. refコールバック関数のクリーンアップ
- refコールバック関数からクリーンアップ関数を返すことで、要素がアンマウントされたときにリソースを解放
- ただし、refコールバック関数は、要素がマウント/アンマウントされたときに呼び出されるため、クリーンアップのタイミングが適切

### 2. SSRの考慮
- localStorageやwindowオブジェクトへのアクセスは、SSR時にはエラーになる可能性がある
- `typeof window !== 'undefined'`でチェックする

### 3. 状態の初期化
- refコールバック関数を使用する場合、状態の初期化タイミングに注意
- 必要に応じて、`useRef`を使用して初期化済みフラグを管理

### 4. パフォーマンス
- refコールバック関数は、要素がマウント/アンマウントされたときに呼び出されるため、不要な再レンダリングを避ける
- ただし、refコールバック関数内で状態を更新する場合は、再レンダリングが発生する

### 5. テスト
- refコールバック関数を使用する場合、テストコードも更新が必要
- テストコードでも、useEffectを排除する（絶対に必要な場合を除く）

## テスト戦略

### 1. 単体テスト
- 各コンポーネントとカスタムフックの動作を確認
- refコールバック関数が正しく呼び出されることを確認
- クリーンアップが正しく実行されることを確認

### 2. 統合テスト
- ページ全体の動作を確認
- データ読み込みが正しく実行されることを確認
- イベントハンドラーが正しく動作することを確認

### 3. E2Eテスト
- 既存のE2Eテストが正常に動作することを確認
- 必要に応じて、E2Eテストを更新

## 移行計画

### Phase 1: 初期データ読み込みの置き換え
1. `client/app/dm-posts/page.tsx`
2. `client/app/dm-users/page.tsx`
3. `client/app/dm-user-posts/page.tsx`
4. `client/app/dm_feed/[userId]/page.tsx`
5. `client/app/dm_feed/[userId]/[postId]/page.tsx`

### Phase 2: Intersection Observerの置き換え
1. `client/app/dm_feed/[userId]/page.tsx`
2. `client/app/dm_feed/[userId]/[postId]/page.tsx`
3. `client/lib/hooks/use-intersection-observer.ts`

### Phase 3: イベントリスナーの置き換え
1. `client/lib/hooks/use-scroll.ts`
2. `client/lib/hooks/use-media-query.ts`

### Phase 4: localStorage読み込みの置き換え
1. `client/lib/hooks/use-local-storage.ts`

### Phase 5: その他のuseEffectの置き換え
1. `client/components/shared/counting-numbers.tsx`
2. `client/app/dm_movie/upload/page.tsx`

### Phase 6: テストコードの更新
1. テストコード内のuseEffectを排除（絶対に必要な場合を除く）

## リスクと対策

### リスク1: refコールバック関数のタイミング問題
**対策**: 必要に応じて、`useRef`を使用して初期化済みフラグを管理

### リスク2: SSR時のエラー
**対策**: `typeof window !== 'undefined'`でチェック

### リスク3: クリーンアップの漏れ
**対策**: refコールバック関数からクリーンアップ関数を返す

### リスク4: 既存のテストの失敗
**対策**: テストコードも更新し、useEffectを排除

## 参考情報

### 関連ドキュメント
- Next.js App Routerドキュメント
- React Hooksドキュメント
- React refコールバック関数ドキュメント

### 関連Issue
- https://github.com/taku-o/go-webdb-template/issues/149: 本設計書の元となったIssue

### 技術スタック
- **フレームワーク**: Next.js 14+ (App Router)
- **言語**: TypeScript 5+
- **UIライブラリ**: shadcn/ui
- **スタイリング**: Tailwind CSS
