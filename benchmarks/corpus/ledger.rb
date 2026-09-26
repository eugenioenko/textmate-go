# frozen_string_literal: true

require "bigdecimal"
require "json"

module Accounting
  class InsufficientFunds < StandardError
    def initialize(account, amount)
      super("#{account.name} cannot cover #{format('%.2f', amount)}")
    end
  end

  Entry = Struct.new(:amount, :memo, :at, keyword_init: true) do
    def credit? = amount.positive?
  end

  class Account
    attr_reader :name, :entries

    def initialize(name, opening: 0)
      @name = name
      @entries = []
      deposit(opening, memo: "opening balance") unless opening.zero?
    end

    def balance
      entries.sum(BigDecimal("0")) { |e| e.amount }
    end

    def deposit(amount, memo: nil)
      raise ArgumentError, "amount must be positive" unless amount.positive?

      record(BigDecimal(amount.to_s), memo)
    end

    def withdraw(amount, memo: nil)
      amount = BigDecimal(amount.to_s)
      raise InsufficientFunds.new(self, amount) if amount > balance

      record(-amount, memo)
    end

    def transfer(to:, amount:)
      withdraw(amount, memo: "transfer to #{to.name}")
      to.deposit(amount, memo: "transfer from #{name}")
    end

    def statement(since: nil)
      rows = entries.select { |e| since.nil? || e.at >= since }
      rows.map do |e|
        sign = e.credit? ? "+" : "-"
        "#{e.at.strftime('%F')} #{sign}#{e.amount.abs.to_s('F')} #{e.memo}"
      end
    end

    def to_json(*args)
      { name:, balance: balance.to_s("F"), entries: entries.size }.to_json(*args)
    end

    private

    def record(amount, memo)
      entries << Entry.new(amount:, memo:, at: Time.now)
      self
    end
  end
end

if $PROGRAM_NAME == __FILE__
  checking = Accounting::Account.new("checking", opening: 100)
  savings = Accounting::Account.new("savings")
  checking.transfer(to: savings, amount: 42.5)
  puts checking.statement, savings.to_json
end
