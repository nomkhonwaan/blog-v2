import { Input } from '@angular/core';

export abstract class UserComponent {

  @Input()
  userInfo: UserInfo;

}
