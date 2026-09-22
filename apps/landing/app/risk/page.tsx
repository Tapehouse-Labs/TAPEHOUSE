import type { Metadata } from "next";
import { Entry, PlainPage } from "../components/PlainPage";

export const metadata: Metadata = {
  title: "Risk disclosure",
  alternates: { canonical: "/risk" },
  description: "What can go wrong when you borrow against Stock Tokens with Tapehouse.",
};

export default function Risk() {
  return (
    <PlainPage title="Read this before you borrow" lede="Plainly, and in full: what Tapehouse is, and what can go wrong.">
      <Entry heading="What Tapehouse is">
        <p>
          Tapehouse is software: a set of smart contracts and this interface to them. It is not a broker
          or a bank, it holds no licence to provide financial services, and nobody at Tapehouse takes
          custody of your assets. Nothing on this site is investment, legal or tax advice, and nothing
          here is an offer of credit. Today it runs on a test network, where assets have no value.
        </p>
      </Entry>
      <Entry heading="What a Stock Token is, and is not">
        <p>
          A Stock Token is not a share. It is a debt security of the company that issues it, built to
          follow the price of a share or a fund. Holding one gives you no vote and no standing as a
          shareholder, and it can only be redeemed for cash, through the providers the issuer authorises.
        </p>
        <p>
          That has consequences. You depend on the issuer staying solvent and honouring the token, which
          you would not if you owned the share. The token can trade away from the share it follows, most
          of all when the underlying market is shut. Splits and other corporate actions reach you through
          the token’s own adjustment, not through a broker. And issuers do not offer Stock Tokens to
          everyone: United States persons, among others, are excluded.
        </p>
      </Entry>
      <Entry heading="What borrowing can cost you">
        <p>
          Borrowing against a portfolio can lose you part or all of it. If the value of your holdings falls
          far enough, some of them will be sold to repay the loan, and that can happen while the
          underlying market is closed.
        </p>
      </Entry>
      <Entry heading="What can fail">
        <p>
          A published band is an estimate, and it is sometimes wrong; that is why we publish the record.
          Price sources can fail, smart contracts can contain errors, and a network can halt.
        </p>
      </Entry>
      <Entry heading="Where it is available">
        <p>Tapehouse may not be available where you live. Check before you use it.</p>
      </Entry>
    </PlainPage>
  );
}
