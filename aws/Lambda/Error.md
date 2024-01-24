### FUNCTION_ERROR_INIT_FAILURE

Provisioned Concurrency configuration failed to be applied. Reason: FUNCTION_ERROR_INIT_FAILURE

Lambda 本身有错误,在配置Concurrency时, 需要先启动lambda,而lambda在初始化时就出现了错误,所以无法配置Concurrency.