export const environment = {
  production: true,

  url: 'https://www.nomkhonwaan.com',
  title: 'Nomkhonwaan | Trust me I\'m Petdo',

  version: '${VERSION}',
  revision: '${REVISION}',

  auth0: {
    clientId: 'cSMgdzCX59n4TcL7H6RWRUYeRFGqCMbU',
    redirectUri: 'https://www.nomkhonwaan.com/login',
    audience: 'https://www.nomkhonwaan.com',
  },

  graphql: {
    endpoint: '/graphql',
  },
};