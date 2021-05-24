import React from 'react';
import logo from './logo.svg';
import './App.css';

const getSSOUri: () => Promise<string | null> = () => {
  return fetch(
    'http://server.saml.test/v1/login_uri'
  ).then(r => {
    return r.json();
  }).then(r => {
    return r.uri;
  })
  .catch(() => null)
};

const isLoggedIn: () => Promise<boolean> = () => {
  return fetch(
    'http://server.saml.test/v1/is_logged_in',
    {
      credentials: 'include'
    }
  ).then(r => {
    return r.json();
  }).then(r => {
    return r.logged_in || false;
  })
  .catch(() => false)
}

const logout = () => {
  fetch(
    'http://server.saml.test/v1/logout',
    {
      credentials: 'include'
    }
  ).then(r => {
    window.location.replace(window.location.href);
  })
  .catch((e) => {
    console.error(e);
    window.location.replace(window.location.href);
  })
}

const App = () => {

  const [loggedIn, setLoggedIn] = React.useState(false)
  const [loading, setLoading] = React.useState(!loggedIn);
  const [loginURI, setLoginURI] = React.useState('');

  React.useEffect(() => {
    setLoading(true);
    isLoggedIn().then(l => {
      setLoggedIn(l);
      setLoading(false);
    })
  }, []);

  React.useEffect(() => {
    if (!loggedIn && !loading && !loginURI) {
      setLoading(true);
      getSSOUri().then(uri => {
        if (uri) {
          setLoginURI(uri) 
        }
        setLoading(false);
      })
    }
  }, [loggedIn, loading]);

  return (
    <div className="App">
      <div style={{marginBottom:'10px'}}>
        {
          loading
          ?
          <p>Please wait</p>
          :
          <p>{
            loginURI ? 'Click below to login' : (!loggedIn && 'No SP server bro')
          }</p>
        }{
          loggedIn && <p>Click below to logout</p>
        }
        {
          !loading && !!loginURI && (
            <a
              href={loginURI}
            >
              Login
            </a>
          )
        }
        {
          loggedIn && (
            <a
              onClick={() => {
                setLoading(true);
                logout();
              }}
            >
              Logout
            </a>
          )
        }
      </div>
    </div>
  );
}

export default App;
