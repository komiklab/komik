from functools import lru_cache
from pydantic_settings import BaseSettings, SettingsConfigDict

class Settings(BaseSettings):
    komik_endpoint: str
    
    model_config = SettingsConfigDict(env_file=".env.development")

@lru_cache
def get_settings() -> Settings:
    return Settings()