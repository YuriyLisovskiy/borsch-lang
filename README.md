# Мова програмування Борщ <img width="200" src="Images/logo.svg" align="right" />

[![docs](https://img.shields.io/badge/%D0%94%D0%BE%D0%BA%D1%83%D0%BC%D0%B5%D0%BD%D1%82%D0%B0%D1%86%D1%96%D1%8F-%D0%91%D0%BE%D1%80%D1%89-blue)](https://yuriylisovskiy.github.io/borsch-lang/)

Борщ - це мова програмування інтерпретованого типу, яка дозволяє писати
код українською мовою.

## Збірка та встановлення з вихідного коду
### Передумови
Перед збіркою, необхідно встановити усе необхідне для розробки на Go.
Інструкції зі встановлення можна знайти [тут](https://go.dev/doc/install).

### Збірка та встановлення
Зібрати та встановити інтерпретатор, а також встановити стандартну бібліотеку:
* для `*nix` (використовуючи утиліту `make`):
  ```bash
  make install
  ```
* для `Windows`:
  ```
  .\Scripts\install.bat
  ```

## Видалення середовища
Видалити інтерпретатор разом зі стандартною бібліотекою:
* для `*nix`:
  ```bash
  make uninstall
  ```
* для `Windows`:
  ```
  .\Scripts\uninstall.bat
  ```

## Порядок розробки
### Розробка
* Створити issue з описом змін
* Створити гілку `<тип>/<коротка назва змін, або номер issue>`, приклади:
  * новий функціонал: `feature/issue-47`, `feature/add-some-operator`
  * виправлення: `fix/issue-48`, `fix/some-critical-bug`
* Додати зміни
* Створити запит на злиття створеної гілки із `dev`
* Залити зміни в `dev`

### Випуск
* Створити запит на злиття `dev` із `master`
* Залити зміни в `master`
* Створити тег у форматі v*.*.*
