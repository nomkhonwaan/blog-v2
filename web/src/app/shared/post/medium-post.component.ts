import { Component, Input, OnInit } from '@angular/core';

import { Type } from './type.enum';

@Component({
  selector: 'app-medium-post',
  templateUrl: './medium-post.component.html',
  styleUrls: ['./medium-post.component.scss'],
})
export class MediumPostComponent implements OnInit {

  @Input()
  post: Post;

  type: Type = Type.Medium;

  ngOnInit(): void {
  }

}
