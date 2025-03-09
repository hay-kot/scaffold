---
---

# Template Partials

Template partials allow you to create reusable template components that can be included in multiple files within your scaffold project. This feature helps reduce duplication and makes your templates more maintainable.

## How Partials Work

Partials are stored in a special `partials` directory at the root of your scaffold project. The scaffold engine automatically detects this directory and registers all files within it as available partials.

:::v-pre
```
my-scaffold/
├── scaffold.yml
├── partials/
│   ├── header.txt
│   ├── footer.txt
│   └── common/
│       └── sidebar.txt
├── {{ .Project }}/
    └── index.html
```
:::

## Using Partials in Templates

To use a partial in your template files, use the `partial` template function:

:::v-pre
```
{{ partial "header" . }}

<!-- Main content here -->

{{ partial "footer" . }}
```
:::

The `partial` function takes two arguments:
1. The name of the partial (without file extension)
2. The data to pass to the partial (typically the current context using `.`)

## Partial Naming and Organization

Partials can be organized in subdirectories within the `partials` folder. When referencing them, use the path relative to the `partials` directory, without file extensions:

:::v-pre
```
{{ partial "common/sidebar" . }}
```
:::

## Passing Data to Partials

Partials have access to the same data context as the template where they're used. When you pass `.` as the second argument to the `partial` function, the partial receives the full context:

:::v-pre
```
<!-- In your partial file (partials/greeting.txt) -->
Hello, {{ .Computed.Username }}!

<!-- In your template file -->
{{ partial "greeting" . }}
```
:::

## Example

Here's a complete example of using partials:

**partials/header.txt**:
:::v-pre
```
<header>
  <h1>{{ .Project }}</h1>
  <p>{{ .Computed.Description }}</p>
</header>
```
:::

:::v-pre
**{{ .Project }}/index.html**:
:::

:::v-pre
```html
<!DOCTYPE html>
<html>
<head>
  <title>{{ .Project }}</title>
</head>
<body>
  {{ partial "header" . }}
  
  <main>
    <!-- Main content here -->
  </main>
</body>
</html>
```
:::

## Best Practices

- Keep partials focused on a single responsibility
- Use meaningful names that reflect their purpose
- Organize related partials in subdirectories for better maintainability
- Document any special requirements or expected variables in comments at the top of each partial