from pathlib import Path
path = Path('apps/im/ws/websocket/server.go')
with path.open(encoding='utf-8') as f:
    for i,line in enumerate(f,1):
        if '//' in line:
            print(f'{i}: {line.rstrip()}')
