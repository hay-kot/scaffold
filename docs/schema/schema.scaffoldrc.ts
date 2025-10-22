export interface Schema {
  settings?: {
    /**
     * The theme for use in the scaffold. The default is `scaffold`.
     * */
    theme?: "scaffold" | "charm" | "dracula" | "base16" | "catpuccino" | "tokyo-night";
    /**
     * The behavior for running hooks. `prompt` will ask the user if they want to run the hooks ,
     * `always` will always run the hooks, and `never` will never run the hooks. The default is
     * `prompt`.
     * */
    run_hooks?: "always" | "never" | "prompt";
    /**
     * the log level for the scaffold. The default is `warn`.
     * */
    log_level?: "debug" | "info" | "warn" | "error";
    /**
     * the log file path for the scaffold. The default is none. When provided without a `/` prefix
     * the path is relative to scaffoldrc file.
     * */
    log_file?: string;
  };
  /**
   * The aliases section allows you to define key/value pairs as shortcuts for a scaffold path.
   * This is useful to shorten a reference for a specific scaffold.
   * */
  aliases?: {
    [key: string]: string;
  };
  /**
   * The defaults section allows you to define key/value pairs that will be used as defaults for
   * prompts.
   * */
  defaults?: {
    [key: string]: string;
  };
  /**
   * The shorts section allows you to define expandable text snippets. Commonly these would be used
   * to prefix a URL or path.
   * */
  shorts?: {
    [key: string]: string;
  };
  /**
   * The auth sections lets you define authentication matchers for your scaffolds. This is useful for
   * using scaffolds that are stored in a private repository. The configuration supports basic
   * authentication and token authentication. Note that in most cases, you want basic authentication,
   * even us you're using a personal access token.
   * */
  auth?: (AuthEntryBasic | AuthEntryToken)[];
}

type AuthEntryBasic = {
  /**
   * Match is a regular expression  that will be used to match the scaffold path. If the path matches
   * the regular expression, the authentication will be used.
   * */
  match: string;
  /**
   * The basic authentication entry. This is the username and password that will be used to authenticate
   * */
  basic: {
    username: string;
    password: string;
  };
};

type AuthEntryToken = {
  /**
   * Match is a regular expression  that will be used to match the scaffold path. If the path matches
   * the regular expression, the authentication will be used.
   * */
  match: string;
  /**
   * The token is the personal access token that will be used to authenticate the user.
   * */
  token: string;
};
