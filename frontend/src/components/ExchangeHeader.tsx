import { useState, type FormEvent } from "react";
import type { ExchangeDto } from "../types";

type ExchangeForm = {
  base: string;
  target: string;
  amount: string;
};

type HeroExchangeProps = {
  onExchange: (form: ExchangeForm) => Promise<void>;
  exchangeResult: ExchangeDto | null;
  exchangeError: string;
  formatAmountInput: (value: string) => string;
  formatMoney: (value: string | null | undefined) => string;
};

export default function ExchangeHeader({
  onExchange,
  exchangeResult,
  exchangeError,
  formatAmountInput,
  formatMoney,
}: HeroExchangeProps) {
  const [form, setForm] = useState<ExchangeForm>({
    base: "",
    target: "",
    amount: "",
  });

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    await onExchange(form);
  }

  return (
    <header className="exchange-header">
      <div>
        <p className="eyebrow">Live FX toolkit</p>
        <h1>Currency Exchange Console</h1>
        <p className="lead">
          Manage currencies, set rates, and calculate conversions from one
          place.
        </p>
      </div>
      <div className="exchange-card">
        <p className="exchange-title">Quick Exchange</p>
        <form className="stack" onSubmit={handleSubmit}>
          <div className="field-row">
            <input
              placeholder="Base"
              value={form.base}
              maxLength={3}
              onChange={(event) =>
                setForm({ ...form, base: event.target.value })
              }
            />
            <input
              placeholder="Target"
              value={form.target}
              maxLength={3}
              onChange={(event) =>
                setForm({ ...form, target: event.target.value })
              }
            />
          </div>
          <input
            placeholder="Amount"
            value={form.amount}
            inputMode="decimal"
            onChange={(event) =>
              setForm({
                ...form,
                amount: formatAmountInput(event.target.value),
              })
            }
          />
          <button className="primary" type="submit">
            Exchange
          </button>
        </form>
        {exchangeResult && (
          <div className="result">
            <p>Result</p>
            <strong>
              {formatMoney(exchangeResult.convertAmount)}{" "}
              {exchangeResult.exchangeRate?.targetCurrency?.sign ||
                exchangeResult.exchangeRate?.targetCurrency?.code}
            </strong>
          </div>
        )}
        {exchangeError && <p className="error">{exchangeError}</p>}
      </div>
    </header>
  );
}
