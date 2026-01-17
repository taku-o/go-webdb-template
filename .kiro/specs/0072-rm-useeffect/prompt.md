/kiro:spec-requirements "https://github.com/taku-o/go-webdb-template/issues/149
のissueの条件でNext.jsのコードを修正する要件定義書を作成してください。
GitHub CLIは入っています。
cc-sddのfeature名は0072-rm-useeffectとしてください。"
think.

テストコードもそれが絶対必要ないなら、useEffectは使わない。
>**本実装の範囲外**:
>- テストコードの変更（テストコードのuseEffectは対象外）

コメントは不要。
>### 6.6 保持が必要なuseEffect
>- [ ] `client/components/video-player/video-player.tsx`のクリーンアップ用useEffectが保持されている
>- [ ] 保持が必要なuseEffectには、コメントで理由が記載されている

要件定義書を承認します。

/kiro:spec-design 0072-rm-useeffect

/kiro:spec-tasks 0072-rm-useeffect

設計書を承認します。
タスクリストを承認します。

/kiro:spec-impl 0072-rm-useeffect 1.1
/kiro:spec-impl 0072-rm-useeffect 1.2
/kiro:spec-impl 0072-rm-useeffect 1.3

/kiro:spec-impl 0072-rm-useeffect 1.4

これはどんな警告がでてる？
>dm_feed/[userId]/page.tsxからは警告が出ていません（警告は別のファイル[postId]/page.tsxのもの）。

/kiro:spec-impl 0072-rm-useeffect 1.5

いったんcommitしてください。

