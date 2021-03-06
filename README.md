# Мова програмування Борщ <img width="200" src="Misc/logo.svg" align="right" />

Борщ - це мова програмування інтерпретованого типу, яка дозволяє писати код українською мовою.

### Встановлення та налаштування
Зібрати інтерпретатор:
```bash
make build
```
У каталозі `./bin` буде міститися інтерпретатор `borsch`.

Для встановлення стандартної бібліотеки, достатньо скопіювати каталог `Lib` до бажаного розташування
на диску та експортувти шлях до неї:
```bash
export BORSCH_LIB="/usr/local/lib/borsch-lang/Lib"
```

### Документація
На даний момент документації не існує, але в майбутньому вона буде написана, а посилання
буде розміщено в цьому пункті. Присутні деякі пакети стандартної бібліотеки, які я додаю
під час розробки інтерпретатора. Ознайомитися з ними можна [тут](./Lib).

### Автор
* Copyright © 2021 Yuriy Lisovskiy
