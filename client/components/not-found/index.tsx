import * as React from 'react' ;
import Helmet from 'react-helmet';
import { IndexLink } from 'react-router';

export default class NotFound extends React.Component {

  render() {
    return <div>
      <Helmet title='404 Page Not Found' />
      <h2 className="notFound">
      404 Page Not Found</h2>
    <IndexLink to='/' className="Link">go home</IndexLink>
    </div>;
  }

}
