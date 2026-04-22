CREATE TABLE IF NOT EXISTS public.login_sessions
(
    "sessionId" BYTEA PRIMARY KEY NOT NULL,
    "userId" INTEGER UNIQUE NOT NULL,
    "expiresAt" TIMESTAMPTZ NOT NULL,
    "createdAt" TIMESTAMPTZ NOT NULL,
    CONSTRAINT fk_login_sessions_user
        FOREIGN KEY ("userId")
        REFERENCES public.users ("id")
        ON DELETE CASCADE
)