# LazyPcUser

- [LazyPcUser](#lazypcuser)
  - [Краткое описание](#краткое-описание)
  - [Установка](#установка)
  - [Команды](#команды)
  - [Функционал](#Реализация-функционала)

## Краткое описание

Бот для телеграмма, реализующий управление звуком на вашем устройстве через взаимодействие с ботом.

## Установка

Для установки требуется склонировать репозиторий. Обратиться к BotFather, создать бота и записать ключ в файл token.txt лежащий в папке с main.go файлом. Запустить server.go и main.go.

## Команды
- /start - начальная команда выводящая приветсвие и описывающая доступные команды
- /volume *num* - изменить громкость звука, подставив в место *num* свое число в процентах от 1-100
- /whatvolume - узнать какой процент громкости установлен сейчас

# Реализация функционала
- Управление звуком PC
- Просмотр установленной громкости
- ~~Поддержка разных платформ(Windows, Mac os)~~
- ~~Управление таймером сна~~
