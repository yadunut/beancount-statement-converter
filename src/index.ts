import { parse } from "https://deno.land/std@0.113.0/flags/mod.ts";
import { readCSVObjects } from "https://deno.land/x/csv@v0.7.0/mod.ts";

import { format } from "https://deno.land/std@0.113.0/datetime/mod.ts";

const SUPPORTED_BANKS = ["dbs", "revolut"];

type Type = "CARD_PAYMENT" | "REWARD" | "TOPUP";
type State = "COMPLETED" | "PENDING";

type RevolutTransaction = {
  type: Type;
  date: Date;
  description: string;
  amount: number;
  currency: string;
  state: State;
  balance: number;
};

async function main(args: string[]) {
  const parsedArgs = parse(args);
  const bank = parsedArgs.bank;
  if (!SUPPORTED_BANKS.includes(bank)) throw new Error("unsupported bank");

  const file = await Deno.open(parsedArgs.file);
  if (bank == "dbs") {
    await parseDbs(file);
  } else if (bank == "revolut") {
    await parseRevolut(file);
  }
}

main(["--bank", "revolut", "--file", "revolut-statement.csv"]).catch(
  console.error,
);

async function parseDbs(file: Deno.File) {
}

async function parseRevolut(file: Deno.File) {
  const transactions = map(readCSVObjects(file), transformRevolut);
  for await (const tx of transactions) {
    if (tx.state == "PENDING") continue;
    let description = tx.description;
    let account = "Expenses:";
    if (description == "Bus/mrt") {
      description = "Bus";
      account = "Expenses:Transport:Public";
    }
    console.log(
      `${format(tx.date, "yyyy-MM-dd")} ! "${description}" ""
  Assets:Revolut:SGD  ${tx.amount} ${tx.currency}
  ${account}
`,
    );
  }
}

function transformRevolut(row: Record<string, string>): RevolutTransaction {
  return {
    type: row.Type as Type,
    date: new Date(row["Completed Date"]),
    description: row.Description as string,
    amount: parseFloat(row.Amount),
    currency: row.Currency as string,
    state: row.State as State,
    balance: parseFloat(row.Balance),
  };
}

async function* map<T, U>(
  iter: AsyncIterable<T>,
  callback: (t: T) => U,
): AsyncIterable<U> {
  for await (const val of iter) {
    yield callback(val);
  }
}
