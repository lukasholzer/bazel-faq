import {FunctionComponent, render} from 'preact';
import {
  Bash,
  Highlighter,
  htmlRender,
  init,
  process,
  Python,
  registerLanguages,
  TypeScript,
} from 'highlight-ts';

import './CodeBox.css';

interface CodeBoxProps {
  code: string;
  language: string;
}

registerLanguages(Bash, Python, TypeScript);

export const CodeBox: FunctionComponent<CodeBoxProps> = ({code, language}) => {
  // initialize highlighter
  const highlighter: Highlighter<string> = init(htmlRender);

  const {value} = process(highlighter, code, language);

  return (
    <div className="code-box">
      <pre>
        <code className="code-box-code">
          {value.split(/\n/).map(line => (
            <span dangerouslySetInnerHTML={{__html: line}} />
          ))}
        </code>
      </pre>
    </div>
  );
};
