import {FunctionComponent} from 'preact';
import {CodeBox} from './CodeBox';
import {TagList} from './TagList';
import './Error.css';

export const Error: FunctionComponent = () => {
  const code = `'""' is not recognized as an internal or external command,`;

  return (
    <div className="error-box">
      <CodeBox code={code} language="bash"></CodeBox>
      <TagList tags={['Windows', 'Bash']} />
      <blockquote>
        This means that it cannot locate the bash correctly
      </blockquote>
      <p>
        if you enter the command <code className="inline">where bash</code> it
        should provide you:{' '}
        <code className="inline">C:\tools\msys64\usr\bin\bash.exe</code>.
      </p>
      <p>
        If this is not the case try to add the bin folder of your{' '}
        <strong>msys64</strong> installation on top of your system environment
        variables.
      </p>
    </div>
  );
};
