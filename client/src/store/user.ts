import type { AxiosResponse } from 'axios';
import { writable } from 'svelte/store';
import type { User } from '../lib/apis/generated/api'
import apis from '../lib/apis/api'

export const user = writable<User>(null);

export async function getMeAction(): Promise<void> {
  let res: AxiosResponse<User, any>;
  try {
    res = await apis.getMe();
  } catch (err: any) {
    if (err.response && err.response.status === 401) {
      window.location.href = '/api/v1/oauth2/generate/code';

      return
    } else {
      throw err
    }
  }
  user.set(res.data);
}
