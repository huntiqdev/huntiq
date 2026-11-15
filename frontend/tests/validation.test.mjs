import test from 'node:test';
import assert from 'node:assert/strict';

// Common regexes from validation constants
const ETHIO_PHONE_PATTERN = /^(09|07)\d{8}$/;
const EMAIL_PATTERN = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

function validateEthioPhone(raw) {
  const cleaned = raw.replace(/[\s\-().]/g, '');
  let local = cleaned;
  if (local.startsWith('+251')) {
    local = '0' + local.slice(4);
  } else if (local.startsWith('251')) {
    local = '0' + local.slice(3);
  }

  if (!ETHIO_PHONE_PATTERN.test(local)) {
    return { valid: false, normalized: '', network: 'unknown' };
  }
  const normalized = '+251' + local.slice(1);
  const network = local.startsWith('07') ? 'safaricom' : 'telebirr';
  return { valid: true, normalized, network };
}

test('validates Ethio Telecom / Telebirr numbers (09...)', () => {
  const res1 = validateEthioPhone('0911223344');
  assert.equal(res1.valid, true);
  assert.equal(res1.normalized, '+251911223344');
  assert.equal(res1.network, 'telebirr');

  const res2 = validateEthioPhone('+251911223344');
  assert.equal(res2.valid, true);
  assert.equal(res2.normalized, '+251911223344');
  assert.equal(res2.network, 'telebirr');
});

test('validates Safaricom Ethiopia numbers (07...)', () => {
  const res1 = validateEthioPhone('0712345678');
  assert.equal(res1.valid, true);
  assert.equal(res1.normalized, '+251712345678');
  assert.equal(res1.network, 'safaricom');

  const res2 = validateEthioPhone('+251 71 234 5678');
  assert.equal(res2.valid, true);
  assert.equal(res2.normalized, '+251712345678');
  assert.equal(res2.network, 'safaricom');
});

test('rejects invalid phone numbers', () => {
  assert.equal(validateEthioPhone('12345').valid, false);
  assert.equal(validateEthioPhone('0811223344').valid, false);
  assert.equal(validateEthioPhone('').valid, false);
});

test('validates email addresses', () => {
  assert.equal(EMAIL_PATTERN.test('user@example.com'), true);
  assert.equal(EMAIL_PATTERN.test('user.name+tag@sub.domain.org'), true);
  assert.equal(EMAIL_PATTERN.test('invalid-email'), false);
  assert.equal(EMAIL_PATTERN.test('@missinguser.com'), false);
});

