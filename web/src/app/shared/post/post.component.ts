import { Input, OnInit, HostListener } from '@angular/core';

export abstract class PostComponent {

  @Input()
  post: Post;

  /**
   * A window inner width number
   */
  innerWidth: number;

  /**
   * A window inner height number
   */
  innerHeight: number;


  constructor() {
    this.onResizeWindow();
  }

  @HostListener('window:resize', [])
  onResizeWindow() {
    this.innerWidth = window.innerWidth;
    this.innerHeight = window.innerHeight;
  }

}
