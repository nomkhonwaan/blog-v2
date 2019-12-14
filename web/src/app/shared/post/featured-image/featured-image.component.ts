import { ChangeDetectionStrategy, Component, HostBinding, OnInit } from '@angular/core';
import { PostTitleComponent } from '../title';

@Component({
  selector: 'app-post-featured-image',
  changeDetection: ChangeDetectionStrategy.OnPush,
  templateUrl: './featured-image.component.html',
  styleUrls: ['./featured-image.component.scss'],
})
export class PostFeaturedImageComponent extends PostTitleComponent implements OnInit {

  /**
   * A featured image source URL
   */
  src: string;

  /**
   * Use to binding to host class names when featured image exists
   */
  @HostBinding('class.-with-featured-image')
  withFeaturedImage = false;


  ngOnInit(): void {
    super.ngOnInit();

    if (this.hasFeaturedImage()) {
      this.src = `/api/v2.1/storage/${this.post.featuredImage.slug}?width=${this.windowInnerWidth}`;

      this.withFeaturedImage = true;
    }
  }

  hasFeaturedImage(): boolean {
    return this.post.featuredImage.slug.length > 0;
  }

}
