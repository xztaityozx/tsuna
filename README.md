# tsuna

`tsuna`、それはデーモンスレイヤーワタナベ

## 注意！
プロセスをエイヤで次々止めるジョークツールです。使うのは推奨しません

## Usage

```zsh
# 何も指定しないときランダムにワタナベさんを生成します
$ tsuna 
私は 渡邉 HOGE (ワタナベ ホゲ) です

# higeサブコマンドは指定したPidのプロセスにいくつかのシグナルを送ります
# QUIT => HUP => INT => TERM => KILL の順番で送ります
$ tsuna hige 99999

# watchサブコマンドはプロセス一覧からプロセス名が一致するものを見つけ次第いくつかのシグナルを送ります
# QUIT => HUP => INT => TERM => KILL の順番で送ります
$ tsuna watch sleep
```

それぞれのオプションについては `--help` をご覧ください

## Install

```zsh
$ go install github.com/xztaityozx/tsuna
```

## License
[MIT](./LICENSE)
