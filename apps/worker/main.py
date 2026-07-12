import asyncio
import logging
from loguru import logger
from komik_worker.komik_worker import KomikWorker

logging.basicConfig(level=logging.INFO)
async def main()->None:
    logger.info("Starting worker")
    worker = KomikWorker()
    await worker.run()


if __name__ == "__main__":
    asyncio.run(main())
