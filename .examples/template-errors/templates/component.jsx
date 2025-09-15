import React from 'react';

// This will cause an error due to JSX syntax
const Component = () => {
  return (
    <div>
      <h1>{{ .Scaffold.project_name }}</h1>
      {/* This JSX expression will conflict with template delimiters */}
      <button onClick={() => {{ console.log('clicked') }}}>
        Click me
      </button>
    </div>
  );
};

export default Component;