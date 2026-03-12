from pathlib import Path
path = Path('apps/im/ws/websocket/server.go')
lines = path.read_text(encoding='utf-8').splitlines()
import itertools
ranges = [(1,80),(80,160),(160,260),(260,360)]
for start,end in ranges:
    print(f'===== {start}-{end} =====')
    for idx in range(start, min(end, len(lines))+1):
        print(f'{idx}: {lines[idx-1]}')
