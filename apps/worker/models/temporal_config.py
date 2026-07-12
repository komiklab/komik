from pydantic import Field
from pydantic_settings import BaseSettings,SettingsConfigDict

class Configuration(BaseSettings):
    model_config = SettingsConfigDict(
        env_prefix="TEMPORAL_",
        env_file=".env.development",
        env_file_encoding="utf-8",
        env_nested_max_split=1,
    )

    url: str = Field(
        default="localhost:7233",
        description="Target for the Temporal client"
    )

    namespace: str = Field(
        default="default",
        description="Namespace for the Temporal client"
    )
    
    task_queue: str = Field(
        default="komik-TaskQueue",
        description="Task queue for the Temporal client"
    )