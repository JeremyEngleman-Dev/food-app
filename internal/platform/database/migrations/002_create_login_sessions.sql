CREATE TABLE IF NOT EXISTS public.login_sessions
(
    "sessionId" TEXT PRIMARY KEY,
    "userId" INTEGER UNIQUE NOT NULL,
    "expiresAt" TIMESTAMPTZ NOT NULL,
    "createdAt" TIMESTAMPTZ NOT NULL,
    CONSTRAINT fk_login_sessions_user
        FOREIGN KEY ("userId")
        REFERENCES public.users ("id")
        ON DELETE CASCADE,
    CONSTRAINT login_unique_session UNIQUE ("sessionId"),
    CONSTRAINT login_unique_user UNIQUE ("userId")
)