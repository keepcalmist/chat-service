package middlewares

// FIXME: В мидлваре, логирующей запрос, необходимо заиспользовать internal/errors.GetServerErrorCode,
// FIXME: чтобы при наличии ошибки менять status на соответствующий код.
// FIXME: Иначе в логах мы всегда будем видеть 200 OK и пропускать ошибки :)
