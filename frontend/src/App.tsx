import { useEffect, useState, type FormEvent } from "react";
import CurrencyPanel from "./components/CurrencyPanel";
import ExchangeHeader from "./components/ExchangeHeader";
import RatePanel from "./components/RatePanel";
import type { CurrencyDto, ExchangeDto, Page } from "./types";

const apiBase = "/api";

type CurrencyForm = {
  code: string;
  fullName: string;
  sign: string;
};

type RateForm = {
  baseCode: string;
  targetCode: string;
  rate: string;
};

type ExchangeForm = {
  base: string;
  target: string;
  amount: string;
};

function formatMoney(value: string | null | undefined) {
  if (value === null || value === undefined || value === "") return "0";
  return value;
}

function normalizeDecimalInput(value: string) {
  return value.replace(",", ".");
}

function formatRateInput(value: string, maxDecimals = 6) {
  const normalized = normalizeDecimalInput(value);
  const parts = normalized.split(".");
  if (parts.length === 1) {
    return normalized;
  }
  const head = parts.shift() || "";
  const fraction = parts.join("");
  return `${head}.${fraction.slice(0, maxDecimals)}`;
}

export default function App() {
  const [currencies, setCurrencies] = useState<CurrencyDto[]>([]);
  const [currencyError, setCurrencyError] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [currencyMessage, setCurrencyMessage] = useState("");
  const [currencyPageNumber, setCurrencyPageNumber] = useState(1);
  const [currencyPageSize, setCurrencyPageSize] = useState(6);
  const [currencyTotal, setCurrencyTotal] = useState(0);
  const [newCurrency, setNewCurrency] = useState<CurrencyForm>({
    code: "",
    fullName: "",
    sign: "",
  });

  const [newRate, setNewRate] = useState<RateForm>({
    baseCode: "",
    targetCode: "",
    rate: "",
  });
  const [rateMessage, setRateMessage] = useState("");
  const [rateError, setRateError] = useState("");

  const [exchangeResult, setExchangeResult] = useState<ExchangeDto | null>(null);
  const [exchangeError, setExchangeError] = useState("");

  async function loadCurrencies(
    pageNumber = currencyPageNumber,
    pageSize = currencyPageSize
  ) {
    setIsLoading(true);
    setCurrencyError("");
    try {
      const response = await fetch(
        `${apiBase}/currencies?pageNumber=${pageNumber}&pageSize=${pageSize}`
      );
      if (!response.ok) {
        const payload = await response.json();
        throw new Error(payload.message || "Failed to load currencies");
      }
      const payload = (await response.json()) as Page<CurrencyDto>;
      setCurrencies(payload.items || []);
      setCurrencyPageNumber(payload.pageNumber || pageNumber);
      setCurrencyPageSize(payload.pageSize || pageSize);
      setCurrencyTotal(payload.total || 0);
    } catch (error) {
      if (error instanceof Error) {
        setCurrencyError(error.message);
      }
    } finally {
      setIsLoading(false);
    }
  }

  useEffect(() => {
    void loadCurrencies();
  }, []);

  useEffect(() => {
    void loadCurrencies(currencyPageNumber, currencyPageSize);
  }, [currencyPageNumber, currencyPageSize]);

  async function handleCreateCurrency(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setCurrencyMessage("");
    setCurrencyError("");
    try {
      const response = await fetch(`${apiBase}/currencies`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(newCurrency),
      });
      const payload = await response.json();
      if (!response.ok) {
        throw new Error(payload.message || "Failed to create currency");
      }
      setCurrencyMessage(`Added ${payload.code}`);
      setNewCurrency({ code: "", fullName: "", sign: "" });
      await loadCurrencies();
    } catch (error) {
      if (error instanceof Error) {
        setCurrencyError(error.message);
      }
    }
  }

  async function handleCreateRate(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setRateMessage("");
    setRateError("");
    try {
      const response = await fetch(`${apiBase}/rates`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          ...newRate,
          rate: Number(newRate.rate),
        }),
      });
      const payload = await response.json();
      if (!response.ok) {
        throw new Error(payload.message || "Failed to create rate");
      }
      setRateMessage(`Created rate`);
      setNewRate({ baseCode: "", targetCode: "", rate: "" });
    } catch (error) {
      if (error instanceof Error) {
        setRateError(error.message);
      }
    }
  }

  async function handleExchange(form: ExchangeForm) {
    setExchangeError("");
    setExchangeResult(null);
    const params = new URLSearchParams({
      base: form.base,
      target: form.target,
      amount: form.amount,
    });
    try {
      const response = await fetch(`${apiBase}/exchange?${params.toString()}`);
      const payload = await response.json();
      if (!response.ok) {
        throw new Error(payload.message || "Failed to exchange");
      }
      setExchangeResult(payload as ExchangeDto);
    } catch (error) {
      if (error instanceof Error) {
        setExchangeError(error.message);
      }
    }
  }

  return (
    <div className="page">
      <ExchangeHeader
        onExchange={handleExchange}
        exchangeResult={exchangeResult}
        exchangeError={exchangeError}
        formatAmountInput={(value) => formatRateInput(value, 6)}
        formatMoney={formatMoney}
      />
      <main className="grid">
        <CurrencyPanel
          currencies={currencies}
          isLoading={isLoading}
          error={currencyError}
          message={currencyMessage}
          pageNumber={currencyPageNumber}
          pageSize={currencyPageSize}
          total={currencyTotal}
          form={newCurrency}
          onPageChange={setCurrencyPageNumber}
          onPageSizeChange={(size) => {
            setCurrencyPageNumber(1);
            setCurrencyPageSize(size);
          }}
          onRefresh={() => loadCurrencies(currencyPageNumber, currencyPageSize)}
          onChange={setNewCurrency}
          onSubmit={handleCreateCurrency}
        />
        <RatePanel
          form={newRate}
          message={rateMessage}
          error={rateError}
          onChange={setNewRate}
          onSubmit={handleCreateRate}
          formatRateInput={(value) => formatRateInput(value, 6)}
        />
      </main>
    </div>
  );
}
