import { signInAction, signOutAction } from '@/lib/actions/auth-actions'
import { Button } from '@/components/ui/button'

interface AuthButtonsProps {
  user: {
    name?: string | null
    email?: string | null
  } | null
}

export function AuthButtons({ user }: AuthButtonsProps) {
  if (user) {
    return (
      <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
        <div>
          <p className="font-semibold text-sm sm:text-base">
            <span className="sr-only">現在のログイン状態: </span>
            ログイン中: {user.name}
          </p>
          {user.email && (
            <p className="text-xs sm:text-sm text-muted-foreground" aria-label="メールアドレス">
              {user.email}
            </p>
          )}
        </div>
        <form action={signOutAction}>
          <Button type="submit" variant="destructive" size="sm" className="w-full sm:w-auto" aria-label="ログアウト">
            ログアウト
          </Button>
        </form>
      </div>
    )
  }

  return (
    <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
      <p className="text-sm sm:text-base text-muted-foreground">
        <span className="sr-only">現在のログイン状態: </span>
        ログインしていません
      </p>
      <form action={signInAction}>
        <Button type="submit" size="sm" className="w-full sm:w-auto" aria-label="ログイン">
          ログイン
        </Button>
      </form>
    </div>
  )
}
