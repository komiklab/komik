from datetime import datetime
from typing import Any
from uuid import UUID

from pydantic import BaseModel, ConfigDict, field_validator


class WorkflowRequest(BaseModel):
    model_config = ConfigDict(extra="allow")

    id: UUID
    created_at: datetime
    updated_at: datetime
    completed_at: datetime | None

    source_type: str
    source_ref: str
    source_payload: dict[str, Any]

    status: str

    temporal_workflow_id: str | None
    temporal_run_id: str | None
    temporal_task_queue: str

    agent_thread_id: str | None
    active_interrupt_id: str | None

    result_summary: str | None
    result: Any | None
    error_message: str | None

    @field_validator("id", mode="before")
    @classmethod
    def convert_uuid_bytes(cls, value):
        if isinstance(value, list):
            return UUID(bytes=bytes(value))
        return value