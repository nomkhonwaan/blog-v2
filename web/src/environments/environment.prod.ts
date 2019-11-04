export const environment = {
  production: true,

  url: 'https://beta.nomkhonwaan.com',

  version: '${VERSION}',
  revision: '${REVISION}',

  auth0: {
    clientId: 'cSMgdzCX59n4TcL7H6RWRUYeRFGqCMbU',
    redirectUri: 'https://www.nomkhonwaan.com/login',
  },

  graphql: {
    endpoint: '/graphql',
  },
};
