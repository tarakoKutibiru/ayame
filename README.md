# WebRTC Signaling Server Ayame

Ayameをゲームのリアルタイム通信で利用しようと思い、試行錯誤していたところ、そのままの仕様だとどうしても不都合があることに気づいたので、改修していく。

# オリジナルとの変更点

## "register"コマンドの廃止
roomの作成と参加が"register"コマンドで一括して行われていて、1対1のピアは区別することができなかった。ゲームに使う場合は、roomを作成した"Host"とroomに参加する"Client"で明確に区別したほうが、便利そうだったので、従来の"register"コマンドを廃止して、roomの作成コマンド"make_room",既存のroomに参加するコマンド"join_room"に切り分けた。適当に"make"としたけど、Go言語的にはmakeは予約語だったので、あとで"create_room"に変更するかも。