export interface Schema {
  /**
   * questions are the prompts that will be used to gather information from the user before generating the project
   * */
  questions: Question[];
  /**
   * computed is a map of values that will be computed based on the answers provided by the user
   * */
  computed?: {
    [key: string]: string;
  };
  rewrites?: Rewrite[];
  /**
   * skips is an array of globs that will be used to skip template rendering for files that match the glob.
   * */
  skips?: string[];
  inject?: Inject[];
  messages?: {
    /**
     * pre is a message that will be displayed before the questions are asked. This field supports markdown syntax
     * */
    pre?: string;
    /**
     * post is a message that will be displayed after the template has been generated. This field supports markdown syntax
     * */
    post?: string;
  };
  features?: Feature[];
  /**
   * presets is a map of key/value pairs that can be used to provide default values for the questions. Generally, this is only used for testing scaffolds, but could be used for other purposes.
   * */
  presents?: {
    [key: string]: {
      [key: string]: string;
    };
  };
}

type Question = {
  /**
   * name is the key that will be used to store the answer provided by the user. This key will be used to reference the answers in the template.
   * */
  name: string;
  /**
   * prompt is the type of prompt that will be used to ask the user for input.
   * */
  prompt: Prompt;
  /**
   * @deprecated use validate.required instead
   * */
  required?: boolean;
  /**
   * when is a conditiona template that is evaluated to determine if the question should be asked.
   * */
  when?: string;
  /**
   * group is an optional key that can be used to group questions together in a shared view.
   * */
  group?: string;
  validate?: Validator;
};

type Validator = {
  /**
   * required (when true) will ensure that the user provides a value for the prompt. Only some prompt types support this option.
   * */
  required?: boolean;
  match?: {
    /**
     * regex is a regular expression that the user's input must match.
     */
    regex: string;
    /**
     * message is a message that will be displayed to the user if the input does not match the pattern. Only some prompt types support this option.
     * */
    message?: string;
  };
  /**
   * min is the minimum value that the user's input must be greater than or equal to.
   *
   * - For strings, this is the minimum length of the string.
   * - For multi-selects, this is the minimum number of items that must be selected.
   * - For looped prompts, this is the minimum number of items that must be provided.
   */
  min?: number;
  /**
   * max is the maximum value that the user's input must be less than or equal to.
   *
   * - For strings, this is the maximum length of the string.
   * - For multi-selects, this is the maximum number of items that can be selected.
   * - For looped prompts, this is the maximum number of items that can be provided.
   */
  max?: number;
};

type Prompt =
  | InputText
  | InputTextMulti
  | InputTextLoop
  | InputConfirm
  | InputSelect
  | InputSelectMulti;

type InputMixinBase = {
  /**
   * message to be displayed to the user, this is the primary message shown to the user. You can think of this as the label for the input/prompt.
   * */
  message: string;
  /**
   * description is an optional message that can be displayed to the user. This is a secondary label that can provide some more context to the user.
   * */
  description?: string;
};

type InputMixinDefault<T> = {
  /**
   * default is the default value that will be used if the user does not provide a value.
   * */
  default?: T;
};

type InputMixinMulti = {
  /**
   * multi (when true) modifies the prompt to allow the user to provide multiple answers. Only some prompt types support this option.
   * */
  multi: true;
};

type InputText = InputMixinBase & InputMixinDefault<string>;

type InputTextMulti = InputText & InputMixinMulti;

type InputTextLoop = InputMixinBase &
  InputMixinDefault<string[]> & {
    /**
     * loop (when true) will keep asking the user for input until they provide an empty string. Only some prompt types support this option.
     * */
    loop: boolean;
  };

type InputConfirm = Omit<InputMixinBase, "message"> &
  InputMixinDefault<boolean> & {
    /**
     * message to be displayed to the user, this is the same as the `message` property, but converts the value into a boolean.
     * */
    confirm: string;
  };

type InputSelect = InputMixinBase &
  InputMixinDefault<string> & {
    multi?: false;
    /**
     * options is an array of strings that will be used to populate the select options.
     * */
    options: string[];
  };

type InputSelectMulti = InputMixinBase &
  InputMixinMulti &
  InputMixinDefault<string[]> & {
    /**
     * options is an array of strings that will be used to populate the select options.
     * */
    options: string[];
  };

type Rewrite = {
  /**
   * The path to the template file
   * */
  from: string;
  /**
   * a path to the destination file, this field supports template syntax
   * */
  to: string;
};

type Inject = {
  /**
   * The mode to use when injecting the code. This can be one of the following:
   * */
  mode?: "before" | "after";
  name: string;
  /**
   * The relative path to the file to inject into from the output directory
   * */
  path: string;
  /**
   * The location to inject the code/text. This is evaluated using the strings.Contains function. Note that ALL matches will be replaced.
   * */
  at: string;
  /**
   * The code/text to inject into the file
   * */
  template: string;
};

type Feature = {
  /*
   * A template value that when evaluates to true will enable the feature
   * */
  value: string;
  /**
   * a list of globs that are considered to be apart of this feature. When the `value` evaluates to true, files matching these globs will be included in the output.
   * */
  globs: string[];
};
