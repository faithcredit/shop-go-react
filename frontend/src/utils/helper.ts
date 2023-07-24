const CURRENCRY_FORMATTER = new Intl.NumberFormat(undefined, {
  currency: 'USD',
  style: 'currency',
});

export const formatCurrencry = (number: any) => {
  return CURRENCRY_FORMATTER.format(number);
};

export const getDate = (date: Date) => {
  return new Date(date).toLocaleDateString('en');
};

export const baseUrl = `${process.env.REACT_APP_SERVER_ENDPOINT}/api`
//   import.meta.env.VITE_MODE === 'development'
//     ? 'http://localhost:5000'
//     : 'https://typeshop-server.onrender.com';
