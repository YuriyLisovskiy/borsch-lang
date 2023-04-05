# Вступ

## Командний рядок
TODO

## Перша програма

Розпочнемо з короткої програми, яка друкує рядок "Привіт, Світе!".
Для виведення тексту, використаємо функцію `друкр`, яка приймає
один параметр будь-якого типу, друкує його в консоль, та додає
спеціальний символ, що позначає новий рядок - `\n`.

Створимо файл `програма.борщ` з потрібним кодом:
```shell
echo 'друкр("Привіт, Світе!");' > програма.борщ 
```

Перевіримо вмістиме створеного файлу:
=== "MacOS / Linux"
    ```shell
    cat програма.борщ
    ```
=== "Windows"
    ```shell
    type програма.борщ
    ```

Отримаємо наступний вивід:
```text
друкр("Привіт, Світе!");
```

Запустимо програму з файлу:
```shell
borsch run -f програма.борщ
```

У результаті виконання вищенаведеного коду, отримаємо наступний результат в терміналі:
```text
Привіт, Світе!
```