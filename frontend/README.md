# Frame Front-end Library

The **Frame Front-end Library** is the App Nerds JavaScript component and CSS library. 
It contains reusable JavaScript components and styles for building website and applications. 
All code is pure ES6 JavaScript and CSS3. It is designed to work well with applications 
using the Go Frame framework.

## 🚀 Getting Started

If you want to use the "kitchen sink" approach, simply start by including the base CSS
and **frame.min.js** in your HTML like below (updating your pathing to match yours set up).

```html
<!DOCTYPE html>
<html lang="en">

<head>
	<meta charset="UTF-8" />
	<meta name="viewport" content="width=device-width, initial-scale=1.0" />
	<title>Base Example</title>
	<link rel="stylesheet" href="/static/base.min.css" />
	<link rel="stylesheet" href="/static/icons.min.css" />
	<link rel="stylesheet" href="/static/components.min.css" />
</head>

<body>
	<main>
		<h1>Base Example</h1>
  </main>

  <script type="module">
		import { Alerter } from "/static/frame.min.js";

		document.addEventListener("DOMContentLoaded", () => {
			const alerter = new Alerter();

      alerter.info("Hello!");
    });

  </script>
</body>
</html>
```

## CSS

Frame includes a pretty basic CSS starting template in **base.css** which covers
a reset and styles for most of your basic elements, such as headings, tables, 
and buttons. It makes use of CSS variables to drive the the colors and look
of these elements. Most base elements do not require any additional classes
to achieve their look.

### Headings
Headings start with a font size of *2.2rem* for H1 and move down from there to 
*1.4rem* for H3. H1 have a color of `--primary-color`.

### Buttons
By default all buttons have a slightly larger font size (*1.1rem*) and a slightly
rounded border of `--border-color`. They also have a click ripple effect. I've also
defined three custom classes for common button uses.

* `.action-button` - Used to denote a primary action. Background is `--action-color`, 
  foreground is `--white`.
* `.success-button` - Used to denote a general success. Background is `--success-color`,
  foreground is `--white`.
* `.danger-button` - Used to denote a dangerious operation. Background is `--danger-color`,
  foreground is `--white`.

### Tables
Tables are styled out of the box. A table should be constructed like so.

```html
<table>
  <caption>Table Example</caption>
  <thead>
    <tr>
      <th scope="col">Col 1</th>
      <th scope="col">Col 2</th>
      <th scope="col">Col 3</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <th scope="row">Value 1</th>
      <td>Value 2</td>
      <td>Value 3</td>
    </tr>
    <tr>
      <th scope="row">Another set</th>
      <td>Value 2</td>
      <td>Value 3</td>
    </tr>
  </tbody>
</table>
```

Table styles are driven by the following variables.

* `--table-border-color` - The border of the table and the nth row color
* `--table-header-background-color` - The background color of the header row
* `--table-header-color` - The foreground color of the header row

### Forms
Forms are designed to work with very little markup out of the box. Here is an 
example form.

```html
<form>
  <label for="first">First</label>
  <input type="text" id="first" autocomplete="name" />
  <small>
    This is additional information for this form field
  </small>

  <label for="second">Second</label>
  <select id="second" autocomplete="on">
    <option value="value1">Value 1</option>
    <option value="value2">Value 2</option>
    <option value="value3">Value 3</option>
  </select>

  <label for="third">Third</label>
  <textarea id="third" autocomplete="off">Text goes here</textarea>

  <footer>
    <button type="button"><i class="icon--mdi icon--mdi--arrow-left" alt="Left arrow"></i> Cancel</button>
    <button class="action-button">Action!</button>
  </footer>
</form>
```

Form elements are designed to be legible with a little larger font and
extra padding. Form elements have a focus border and box shadow to clearly
indicate where the current focus is. Finally, the actions to take on a form
should be placed in a footer element.

### Cards
Frame provides a basic concept of a **Card** element. It is designed to have
a title, body content, and a footer. Here is an example.

```html
<div class="cards">
  <article>
    <h3>Card Demonstration</h3>
    <p>
      This is the first card. It contains content.
      The markup is simple.
    </p>

    <footer>
      <button>Secondary</button>
      <button class="action-button">Primary</button>
    </footer>
  </article>
</div>
```

The first thing to do is wrap all of your cards in a div with a class of `.cards`.
Then, for each card, use an `<article>` tag. Each card should have a title using
a heading. This will likely be a **H3** or lower (never an H1 for sure). The body
can be whatever content you want. Finally, add a footer with buttons to take
actions if you wish.

## Components
Components are either custom web components or Javascript libraries to enhance
your application.
