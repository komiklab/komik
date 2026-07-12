from temporalio import workflow, activity
from temporalio.worker import (
    Interceptor, 
    WorkflowInboundInterceptor, 
    ActivityInboundInterceptor,
    ExecuteActivityInput,
    ExecuteWorkflowInput
)
from shared.logging import _current_logger, _logger

class LoguruWorkflowInboundInterceptor(WorkflowInboundInterceptor):
    async def execute_workflow(self, input: ExecuteWorkflowInput):
        info = workflow.info()
        bound_logger = _logger.bind(
            workflow_id=info.workflow_id,
            run_id=info.run_id,
            workflow_type=info.workflow_type,
        )
        token = _current_logger.set(bound_logger)
        try:
            return await super().execute_workflow(input)
        finally:
            _current_logger.reset(token)

class LoguruActivityInboundInterceptor(ActivityInboundInterceptor):
    async def execute_activity(self, input: ExecuteActivityInput):
        info = activity.info()
        bound_logger = _logger.bind(
            activity_id=info.activity_id,
            workflow_id=info.workflow_id,
            workflow_run_id=info.workflow_run_id,
            activity_type=info.activity_type,
        )
        token = _current_logger.set(bound_logger)
        try:
            return await super().execute_activity(input)
        finally:
            _current_logger.reset(token)

class LoguruInterceptor(Interceptor):
    def workflow_interceptor_class(self, input):
        return LoguruWorkflowInboundInterceptor

    def intercept_activity(self, next: ActivityInboundInterceptor) -> ActivityInboundInterceptor:
        return LoguruActivityInboundInterceptor(next)
