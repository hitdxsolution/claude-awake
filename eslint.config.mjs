// ESLint flat config — geumdo-erp 규약을 이 CLI 에 맞춰 가져온 것.
//   핵심: any 금지 · 명시적 반환타입 · floating promise 금지 · catch 는 unknown · prettier 를 규칙으로 강제.
//   이 프로젝트는 React/Nest 가 아니므로 그쪽 전용 플러그인(react-hooks 등)은 뺐다.

import js from '@eslint/js';
import tseslint from 'typescript-eslint';
import prettierPlugin from 'eslint-plugin-prettier';
import prettierConfig from 'eslint-config-prettier';

export default tseslint.config(
  { ignores: ['dist/**', 'node_modules/**'] },
  js.configs.recommended,
  // 타입 정보를 쓰는 규칙(no-unsafe-*, no-floating-promises 등)까지 켠다.
  ...tseslint.configs.strictTypeChecked,
  prettierConfig,
  {
    files: ['**/*.ts', '**/*.mjs'],
    languageOptions: {
      parserOptions: { projectService: true, tsconfigRootDir: import.meta.dirname },
    },
    plugins: { prettier: prettierPlugin },
    rules: {
      // 포매팅도 린트 에러로 — .prettierrc(singleQuote·trailingComma:all·tabWidth:2·printWidth:150) 를 따른다.
      'prettier/prettier': 'error',

      // any 완전 금지(geumdo-erp 전역 규칙).
      '@typescript-eslint/no-explicit-any': 'error',

      // 공개 함수는 반환타입을 명시한다.
      '@typescript-eslint/explicit-function-return-type': ['error', { allowExpressions: true, allowTypedFunctionExpressions: true }],

      // 비동기 결과는 await 하거나 void 로 명시.
      '@typescript-eslint/no-floating-promises': 'error',

      // catch 는 unknown 으로 받는다.
      '@typescript-eslint/use-unknown-in-catch-callback-variable': 'error',

      // 타입 전용 import 는 `import type` 으로(verbatimModuleSyntax 와 짝).
      '@typescript-eslint/consistent-type-imports': ['error', { fixStyle: 'inline-type-imports' }],
    },
  },
  {
    // CLI 라 stdout 이 곧 로그 채널이다(LaunchAgent/작업 스케줄러가 파일로 캡처).
    //   그래도 console 은 막고 process.stdout.write 만 쓰도록 유지한다.
    files: ['**/*.ts'],
    rules: { 'no-console': 'error' },
  },
  {
    // 이 설정 파일 자신은 tsconfig 프로그램 밖이라 타입 기반 규칙을 끈다(문법 검사만).
    files: ['eslint.config.mjs'],
    ...tseslint.configs.disableTypeChecked,
  },
);
