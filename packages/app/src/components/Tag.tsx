import {FunctionComponent} from 'preact';
import './Tag.css';

export const Tag: FunctionComponent = ({children}) => {
  return <span className="tag">{children}</span>;
};
