export type CurrencyDto = {
  id: number;
  code: string;
  fullName: string;
  sign: string;
};

export type ExchangeRateDto = {
  id: number;
  baseCurrency: CurrencyDto;
  targetCurrency: CurrencyDto;
  rate: string;
};

export type ExchangeDto = {
  exchangeRate: ExchangeRateDto;
  amount: string;
  convertAmount: string;
};

export type Page<T> = {
  items: T[];
  pageNumber: number;
  pageSize: number;
  total: number;
};
