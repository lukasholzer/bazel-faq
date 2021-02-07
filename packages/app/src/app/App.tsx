import {FunctionComponent} from 'preact';
import {Error} from '../components/Error';
import './App.css';

export const App: FunctionComponent = () => {
  return (
    <div className="app">
      <h1>Bazel FAQ</h1>
      <Error />
    </div>
  );
};
