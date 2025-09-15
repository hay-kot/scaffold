// This is a JavaScript file with proper template usage

const app = {
  name: "{{ .Scaffold.project_name }}",
  config: {
    debug: true
  }
};

export default app;