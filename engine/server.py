#!/usr/bin/env python
from rpc.engine import EngineService
from rpc import handler
from thrift.transport import TSocket
from thrift.transport import TTransport
from thrift.protocol import TBinaryProtocol
from thrift.server import TServer
import sys
sys.path.append('/usr/lib64/python3.4/site-packages')


__HOST = 'engine'
__PORT = 6666


handler = handler.EngineServiceHandler()
processor = EngineService.Processor(handler)
transport = TSocket.TServerSocket(__HOST, __PORT)
tfactory = TTransport.TBufferedTransportFactory()
pfactory = TBinaryProtocol.TBinaryProtocolFactory()
server = TServer.TSimpleServer(processor, transport, tfactory, pfactory)

print("Starting engine in:", __PORT)
server.serve()
