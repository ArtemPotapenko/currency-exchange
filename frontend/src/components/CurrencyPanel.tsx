import type { FormEvent } from "react";
import type { CurrencyDto } from "../types";

type CurrencyForm = {
  code: string;
  fullName: string;
  sign: string;
};

type CurrencyPanelProps = {
  currencies: CurrencyDto[];
  isLoading: boolean;
  error: string;
  message: string;
  pageNumber: number;
  pageSize: number;
  total: number;
  form: CurrencyForm;
  onPageChange: (pageNumber: number) => void;
  onPageSizeChange: (pageSize: number) => void;
  onRefresh: () => void;
  onChange: (nextForm: CurrencyForm) => void;
  onSubmit: (event: FormEvent<HTMLFormElement>) => void;
};

export default function CurrencyPanel({
  currencies,
  isLoading,
  error,
  message,
  pageNumber,
  pageSize,
  total,
  form,
  onPageChange,
  onPageSizeChange,
  onRefresh,
  onChange,
  onSubmit,
}: CurrencyPanelProps) {
  const totalPages = Math.max(1, Math.ceil(total / pageSize));

  return (
    <section className="panel">
      <div className="panel-head">
        <h2>Currencies</h2>
        <button className="ghost" type="button" onClick={onRefresh}>
          Refresh
        </button>
      </div>
      <div className="panel-toolbar">
        <div className="pager">
          <button
            className="ghost"
            type="button"
            onClick={() => onPageChange(pageNumber - 1)}
            disabled={pageNumber <= 1}
          >
            Prev
          </button>
          <span>
            Page {pageNumber} / {totalPages}
          </span>
          <button
            className="ghost"
            type="button"
            onClick={() => onPageChange(pageNumber + 1)}
            disabled={pageNumber >= totalPages}
          >
            Next
          </button>
        </div>
        <label className="page-size">
          <span>Size</span>
          <select
            value={pageSize}
            onChange={(event) => onPageSizeChange(Number(event.target.value))}
          >
            {[5, 10, 20, 50].map((size) => (
              <option key={size} value={size}>
                {size}
              </option>
            ))}
          </select>
        </label>
      </div>
      {isLoading ? (
        <p className="muted">Loading currencies...</p>
      ) : (
        <div className="pill-grid">
          {currencies.map((currency) => (
            <div key={currency.code} className="pill">
              <span>
                {currency.code} {currency.sign}
              </span>
              <small>{currency.fullName}</small>
            </div>
          ))}
        </div>
      )}
      {error && <p className="error">{error}</p>}
      {message && <p className="success">{message}</p>}
      <form className="stack" onSubmit={onSubmit}>
        <input
          placeholder="Code"
          value={form.code}
          onChange={(event) => onChange({ ...form, code: event.target.value })}
        />
        <input
          placeholder="Full name"
          value={form.fullName}
          onChange={(event) =>
            onChange({ ...form, fullName: event.target.value })
          }
        />
        <input
          placeholder="Sign"
          value={form.sign}
          onChange={(event) => onChange({ ...form, sign: event.target.value })}
        />
        <button className="primary" type="submit">
          Add currency
        </button>
      </form>
    </section>
  );
}
