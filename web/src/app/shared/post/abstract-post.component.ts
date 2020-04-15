import { HostListener, Input } from '@angular/core';

export abstract class AbstractPostComponent {

  /**
   * A post object
   */
  @Input()
  post: Post;

  /**
   * Window width for benefit of on-the-fly image resizing
   */
  innerWidth: number;
  windowInnerWidth: number;

  /**
   * Window height for benefit of on-the-fly image resizing
   */
  innerHeight: number;
  windowInnerHeight: number;

  constructor() {
    this.onResizeWindow();
  }

  @HostListener('window:resize', [])
  onResizeWindow() {
    this.windowInnerWidth = window.innerWidth;
    this.windowInnerHeight = window.innerHeight;
    this.innerWidth = window.innerWidth;
    this.innerHeight = window.innerHeight;
  }

}
