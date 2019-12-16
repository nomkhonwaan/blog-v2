import { ChangeDetectionStrategy, Component, HostBinding, Input, OnInit } from '@angular/core';
import { faFacebookF, faTwitter, IconDefinition as BrandIconDefinition } from '@fortawesome/free-brands-svg-icons';
import { environment } from '../../../../environments/environment';
import { AbstractPostComponent } from '../abstract-post.component';

@Component({
  selector: 'app-post-share-to',
  changeDetection: ChangeDetectionStrategy.OnPush,
  templateUrl: './share-to.component.html',
  styleUrls: ['./share-to.component.scss'],
})
export class PostShareToComponent extends AbstractPostComponent implements OnInit {

  /**
   * Direction of the sharing icons
   */
  @Input()
  flow = 'row';

  /**
   * Use to binding to host class names when direction is specified
   */
  @HostBinding('class.-column')
  get classes(): boolean {
    return this.flow === 'column';
  }

  /**
   * Use to sharing to the social network
   */
  url: string = environment.url;

  /**
   * Use to display an engagement number with human readable format
   */
  shareCount: string;

  /**
   * List of FontAwesome icons
   */
  icons: { [name: string]: BrandIconDefinition } = {
    faFacebookF,
    faTwitter
  };

  ngOnInit(): void {
    const shareCount: number = this.post.engagement.shareCount;

    if (shareCount > 0) {
      if (shareCount < 1000) {
        this.shareCount = shareCount.toString();
      } else if (shareCount < 1000000) {
        this.shareCount = `${(shareCount / 1000).toFixed(2)}k`;
      } else {
        this.shareCount = `${(shareCount / 1000000).toFixed(2)}m`;
      }
    }
  }

  getURL(): string {
    const publishedAt: Date = new Date(
      new Date(this.post.publishedAt)
        .toLocaleString('en-US', { timeZone: 'Asia/Bangkok' }),
    );

    return this.url + '/' + [
      publishedAt.getFullYear().toString(),
      (publishedAt.getMonth() + 1).toString(), // a month number start at 0 not 1
      publishedAt.getDate().toString(),
      this.post.slug,
    ].join('/');
  }

  getEncodedURL(): string {
    return encodeURIComponent(this.getURL());
  }
}
