from pathlib import Path
with Path('apps/im/ws/websocket/server.go').open(encoding='utf-8') as f:
    for i,line in enumerate(f,1):
        if i>40:
            break
        print(f'{i}: {line.rstrip()}')
