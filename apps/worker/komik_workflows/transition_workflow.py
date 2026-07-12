from temporalio import workflow
from models.entity_request import WorkflowRequest
from pydantic import ValidationError

with workflow.unsafe.imports_passed_through():
    from shared.logging import logger

@workflow.defn(name="EntityTransitionWorkflow")
class EntityTransitionWorkflow:
    @workflow.run
    async def run(self, entity: dict) -> None:
        try:
            request = WorkflowRequest.model_validate(entity)
            request.status = "dispatched"
        except ValidationError as e:
            logger.error(f"Invalid entity request: {e}")
            raise
        logger.info(f"EntityTransitionWorkflow started for entity {request.id}")