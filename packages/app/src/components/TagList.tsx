import {FunctionComponent} from 'preact';
import {Tag} from './Tag';
import './TagList.css';

interface TagListProps {
  tags: string[];
}

export const TagList: FunctionComponent<TagListProps> = ({tags}) => {
  return (
    <div className="tag-list">
      {tags.map(tag => (
        <Tag>{tag}</Tag>
      ))}
    </div>
  );
};
