from contextvars import ContextVar
import sys
from loguru import logger as _logger

_logger.remove()
_logger.add(
    sys.stderr,
    format="<green>{time:YYYY-MM-DD HH:mm:ss.SSS}</green> | <level>{level: <8}</level> | <cyan>{name}</cyan>:<cyan>{function}</cyan>:<cyan>{line}</cyan> | <cyan>{extra}</cyan> - <level>{message}</level>"
)

_current_logger = ContextVar("logger", default=_logger)


class LoggerProxy:
    def __getattr__(self, item):
        return getattr(_current_logger.get(), item)


logger = LoggerProxy()