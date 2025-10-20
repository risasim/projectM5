import asyncio
import websockets
import json
from web import web

url = "https://local"
socket = None

webInstance = web(url)
